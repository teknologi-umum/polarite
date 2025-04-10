package slogotel_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	slogotel "polarite/platform/slog-otel"
)

func TestOtelHandler(t *testing.T) {
	const testOperationName = "operation-name"

	setupTracer := func() (*tracetest.SpanRecorder, trace.Tracer) {
		spanRecorder := tracetest.NewSpanRecorder()
		traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
		tracer := traceProvider.Tracer("test-tracer")

		return spanRecorder, tracer
	}

	setupLogger := func(opts ...slogotel.OtelHandlerOpt) *bytes.Buffer {
		var buffer bytes.Buffer
		slog.SetDefault(slog.New(slogotel.New(slog.NewJSONHandler(&buffer, nil), opts...)))

		return &buffer
	}

	t.Run("without span", func(t *testing.T) {
		buffer := setupLogger()

		want := map[string]string{
			"level": "INFO",
			"msg":   "without span",
			"a_key": "a_value",
		}

		func() {
			slog.Info("without span", "a_key", "a_value")
		}()

		got := map[string]string{}
		if err := json.Unmarshal([]byte(strings.TrimSuffix(buffer.String(), "\n")), &got); err != nil {
			panic(err)
		}

		for key := range want {
			if got[key] != want[key] {
				t.Errorf("\ngot %q for key %q\nwant %q", got[key], key, want[key])
			}
		}
	})

	t.Run("adds span and trace ids to log", func(t *testing.T) {
		spanRecorder, tracer := setupTracer()
		buffer := setupLogger()

		want := []string{"trace_id", "span_id"}

		func() {
			ctx := context.Background()
			ctx, span := tracer.Start(ctx, testOperationName)
			defer span.End()

			slog.InfoContext(ctx, "adds span and trace ids to log")
		}()

		got := map[string]string{}
		if err := json.Unmarshal([]byte(strings.TrimSuffix(buffer.String(), "\n")), &got); err != nil {
			t.Fatal(err)
		}

		spans := spanRecorder.Ended()
		if len(spans) != 1 {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), 1)
		}

		for _, key := range want {
			if _, ok := got[key]; !ok {
				t.Errorf("\n%q attribute is missing", key)
			}
		}
	})

	t.Run("adds event to span", func(t *testing.T) {
		spanRecorder, tracer := setupTracer()
		_ = setupLogger()

		want := []attribute.KeyValue{{
			Key:   "a_key",
			Value: attribute.StringValue("a_value"),
		}, {
			Key:   "msg",
			Value: attribute.StringValue("adds event to span"),
		}, {
			Key:   "level",
			Value: attribute.StringValue("INFO"),
		}, {
			Key:   "string_slice",
			Value: attribute.StringSliceValue([]string{"value_1", "value_2"}),
		}, {
			Key:   "int_slice",
			Value: attribute.IntSliceValue([]int{1, 2}),
		}, {
			Key:   "int64_slice",
			Value: attribute.Int64SliceValue([]int64{1, 2}),
		}, {
			Key:   "float64_slice",
			Value: attribute.Float64SliceValue([]float64{1.0, 2.0}),
		}, {
			Key:   "bool_slice",
			Value: attribute.BoolSliceValue([]bool{true, false}),
		}, {

			Key:   "group_1.key_1",
			Value: attribute.StringValue("value_1"),
		}, {
			Key:   "group_2.key_2",
			Value: attribute.StringValue("value_2"),
		}, {
			Key:   "err",
			Value: attribute.StringValue("boom"),
		}}

		func() {
			ctx := context.Background()
			ctx, span := tracer.Start(ctx, testOperationName)
			defer span.End()

			group1 := slog.Group("group_1", "key_1", "value_1")
			group2 := slog.Group("group_2", "key_2", "value_2")

			stringSlice := []string{"value_1", "value_2"}
			intSlice := []int{1, 2}
			int64Slice := []int64{1, 2}
			float64Slice := []float64{1.0, 2.0}
			boolSlice := []bool{true, false}

			slog.InfoContext(ctx,
				"adds event to span",
				"string_slice", stringSlice,
				"int_slice", intSlice,
				"int64_slice", int64Slice,
				"float64_slice", float64Slice,
				"bool_slice", boolSlice,
				"a_key", "a_value",
				"err", errors.New("boom"),
				group1,
				group2,
			)
		}()

		spans := spanRecorder.Ended()

		if len(spans) != 1 {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), 1)
		}

		expectedEventName := "log_record"
		if spans[0].Events()[0].Name != expectedEventName {
			t.Errorf("\ngot %q\nwant %q", spans[0].Events()[0].Name, expectedEventName)
		}

		for _, wantAttr := range want {
			found := false

			for _, gotAttr := range spans[0].Events()[0].Attributes {
				if wantAttr.Key == gotAttr.Key &&
					wantAttr.Value == gotAttr.Value {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("\nspan event attribute with key %v and value %v is missing",
					wantAttr.Key, wantAttr.Value)
			}
		}
	})

	t.Run("adds context baggage attributes to log", func(t *testing.T) {
		spanRecorder, tracer := setupTracer()
		buffer := setupLogger()

		want := map[string]string{
			"key1b": "value1b",
			"key2b": "value2b",
		}

		func() {
			m1, _ := baggage.NewMember("key1b", "value1b")
			m2, _ := baggage.NewMember("key2b", "value2b")
			bag, _ := baggage.New(m1, m2)
			ctx := baggage.ContextWithBaggage(context.Background(), bag)

			ctx, span := tracer.Start(ctx, testOperationName)
			defer span.End()

			slog.InfoContext(ctx, "adds context baggage attributes to log")
		}()

		spanRecorder.Ended()

		got := map[string]string{}
		if err := json.Unmarshal([]byte(strings.TrimSuffix(buffer.String(), "\n")), &got); err != nil {
			t.Fatal(err)
		}

		for key := range want {
			if got[key] != want[key] {
				t.Errorf("\ngot %q for key %q\nwant %q", got[key], key, want[key])
			}
		}
	})

	t.Run("does not set span status with non error logs", func(t *testing.T) {
		spanRecorder, tracer := setupTracer()
		_ = setupLogger()

		want := sdktrace.Status{
			Code: codes.Unset,
		}

		func() {
			ctx := context.Background()

			ctx, span := tracer.Start(ctx, testOperationName)
			defer span.End()

			slog.InfoContext(ctx, "sets span status as error with error log")
			slog.DebugContext(ctx, "sets span status as error with error log")
			slog.WarnContext(ctx, "sets span status as error with error log")
		}()

		spans := spanRecorder.Ended()
		for _, span := range spans {
			if span.Status() != want {
				t.Errorf("\ngot %v\nwant %v", span.Status(), want)
			}
		}
	})

	t.Run("sets span status as error with error log", func(t *testing.T) {
		spanRecorder, tracer := setupTracer()
		buffer := setupLogger()

		want := sdktrace.Status{
			Code:        codes.Error,
			Description: "an error",
		}

		func() {
			ctx := context.Background()

			ctx, span := tracer.Start(ctx, testOperationName)
			defer span.End()

			slog.ErrorContext(ctx, "an error")
		}()

		spans := spanRecorder.Ended()
		spans[0].Status()

		got := map[string]string{}
		if err := json.Unmarshal([]byte(strings.TrimSuffix(buffer.String(), "\n")), &got); err != nil {
			t.Fatal(err)
		}

		if spans[0].Status() != want {
			t.Errorf("\ngot %v\nwant %v", spans[0].Status(), want)
		}
	})

	t.Run("when configured with NoTraceEvents, does not attach events to active trace", func(t *testing.T) {
		wantMsg := "a log without any trace events"
		spanRecorder, tracer := setupTracer()
		buffer := setupLogger(slogotel.WithNoTraceEvents(true))

		func() {
			ctx, span := tracer.Start(context.Background(), testOperationName)
			defer span.End()
			slog.InfoContext(ctx, wantMsg)
		}()

		spans := spanRecorder.Ended()
		spans[0].Status()

		if eventsLen := len(spans[0].Events()); eventsLen > 0 {
			t.Errorf("Expected no events on the span, but there are %d: %v", eventsLen, spans[0].Events())
		}
		got := map[string]string{}
		if err := json.Unmarshal([]byte(strings.TrimSuffix(buffer.String(), "\n")), &got); err != nil {
			t.Fatal(err)
		}
		if msg := got["msg"]; msg != wantMsg {
			t.Errorf("\ngot %v\nwant %v", msg, wantMsg)
		}
	})
}
