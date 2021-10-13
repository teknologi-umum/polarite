package resources_test

import (
	"polarite/resources"
	"testing"
)

var LoremIpsum = `Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec,
pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate
 eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium.
 Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula,
 porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus.
 Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue.
 Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus,
 sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit
 id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante.
 Etiam sit amet orci eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales sagittis
 magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc, quis gravida magna mi a libero. Fusce vulputate
 eleifend sapien. Vestibulum purus quam, scelerisque ut, mollis sed, nonummy id, metus. Nullam accumsan lorem in dui.
 Cras ultricies mi eu turpis hendrerit fringilla. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices
 posuere cubilia Curae; In ac dui quis mi consectetuer lacinia. Nam pretium turpis et arcu. Duis arcu tortor, suscipit
 eget, imperdiet nec, imperdiet iaculis, ipsum. Sed aliquam ultrices mauris. Integer ante arcu, accumsan a, consectetuer
 eget, posuere ut, mauris. Praesent adipiscing. Phasellus ullamcorper ipsum rutrum nunc. Nunc nonummy metus. Vestibulum
 volutpat pretium libero. Cras id dui.`

func TestCompressContent(t *testing.T) {
	c, err := resources.CompressContent([]byte(LoremIpsum))
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if len(c) > len([]byte(LoremIpsum)) {
		t.Error("compressed item is larger than the original. life is meaningless.")
	}
}

func TestDecompressContent(t *testing.T) {
	b := []byte{120, 156, 84, 85, 219, 202, 228, 54, 12, 190, 207, 83, 232, 1, 76, 94, 160, 87, 203, 110, 11, 11, 237, 178, 165, 176, 247, 26, 91, 127, 126, 21, 31, 50, 150, 20, 232, 219, 23, 31, 146, 204, 220, 77, 6, 75, 150, 190, 147, 255, 44, 149, 18, 240, 46, 150, 32, 148, 88, 42, 8, 43, 96, 34, 117, 224, 75, 22, 242, 74, 106, 84, 1, 3, 239, 44, 158, 243, 6, 20, 89, 87, 248, 66, 153, 48, 131, 47, 41, 149, 80, 32, 242, 102, 17, 129, 54, 210, 209, 232, 58, 145, 80, 4, 215, 229, 171, 37, 144, 226, 153, 5, 50, 106, 121, 26, 193, 78, 25, 149, 31, 38, 64, 10, 9, 183, 204, 2, 129, 5, 118, 172, 106, 149, 41, 43, 164, 146, 149, 196, 65, 70, 241, 164, 86, 161, 114, 96, 111, 209, 4, 146, 201, 10, 223, 74, 38, 15, 79, 195, 4, 31, 20, 89, 28, 88, 212, 202, 158, 73, 32, 147, 119, 203, 78, 49, 82, 235, 209, 46, 36, 115, 176, 87, 82, 182, 4, 79, 107, 167, 133, 210, 10, 63, 44, 70, 28, 235, 62, 13, 117, 76, 220, 15, 0, 101, 78, 231, 45, 59, 5, 130, 127, 77, 180, 56, 248, 168, 156, 55, 110, 101, 7, 69, 7, 24, 249, 105, 164, 253, 74, 56, 44, 238, 166, 168, 180, 116, 56, 28, 96, 245, 182, 194, 247, 220, 187, 157, 29, 234, 103, 201, 222, 4, 76, 29, 112, 218, 169, 6, 38, 5, 116, 112, 80, 238, 184, 8, 28, 172, 72, 110, 20, 204, 41, 19, 4, 246, 106, 115, 91, 32, 27, 83, 165, 18, 219, 231, 220, 109, 93, 224, 123, 86, 218, 168, 130, 114, 246, 28, 44, 235, 10, 95, 43, 10, 4, 220, 27, 224, 43, 252, 226, 3, 83, 67, 62, 82, 162, 220, 58, 10, 181, 41, 32, 179, 240, 197, 221, 181, 74, 59, 199, 31, 148, 3, 40, 197, 216, 26, 204, 19, 145, 78, 238, 221, 2, 123, 169, 170, 172, 165, 118, 164, 111, 64, 231, 34, 87, 15, 244, 110, 34, 251, 165, 1, 135, 9, 98, 215, 33, 102, 37, 119, 206, 8, 156, 29, 28, 124, 80, 173, 56, 217, 250, 32, 219, 24, 59, 74, 115, 140, 5, 126, 126, 162, 244, 223, 215, 217, 220, 233, 52, 133, 68, 218, 254, 198, 202, 38, 16, 177, 84, 34, 93, 225, 111, 227, 46, 134, 106, 90, 45, 93, 139, 92, 28, 172, 240, 187, 50, 166, 87, 33, 177, 112, 227, 25, 208, 54, 163, 117, 129, 175, 86, 241, 193, 77, 142, 157, 19, 95, 106, 67, 238, 189, 96, 133, 31, 152, 166, 31, 140, 207, 166, 147, 245, 21, 254, 66, 242, 148, 81, 64, 41, 237, 38, 231, 66, 163, 192, 151, 28, 120, 210, 50, 43, 220, 210, 8, 26, 74, 159, 76, 69, 126, 80, 45, 238, 114, 236, 171, 71, 219, 217, 76, 109, 77, 161, 48, 252, 61, 6, 234, 13, 178, 101, 239, 224, 17, 49, 7, 214, 161, 224, 104, 190, 129, 181, 91, 60, 56, 99, 117, 240, 73, 57, 84, 170, 172, 11, 112, 112, 131, 160, 151, 177, 155, 33, 74, 224, 210, 172, 219, 104, 187, 149, 54, 55, 58, 93, 211, 185, 7, 193, 157, 41, 55, 82, 198, 212, 47, 42, 255, 64, 243, 67, 147, 83, 225, 221, 120, 173, 231, 186, 76, 212, 174, 13, 75, 245, 60, 32, 162, 90, 238, 210, 87, 153, 127, 107, 213, 145, 202, 10, 255, 80, 120, 241, 105, 66, 171, 44, 119, 171, 204, 143, 207, 115, 70, 41, 1, 35, 9, 8, 110, 172, 202, 178, 244, 52, 194, 209, 226, 146, 177, 235, 98, 239, 183, 63, 248, 65, 57, 244, 68, 235, 149, 110, 72, 163, 65, 201, 10, 222, 170, 152, 76, 148, 251, 54, 91, 197, 131, 3, 142, 182, 144, 24, 112, 226, 176, 194, 31, 38, 158, 222, 50, 227, 116, 201, 192, 108, 133, 95, 36, 45, 37, 163, 37, 216, 173, 154, 116, 14, 29, 136, 167, 72, 117, 104, 185, 133, 200, 204, 0, 161, 224, 32, 151, 108, 41, 253, 215, 153, 235, 46, 184, 192, 69, 239, 45, 73, 51, 238, 72, 254, 220, 213, 185, 140, 116, 184, 21, 156, 184, 101, 139, 90, 221, 89, 110, 41, 220, 112, 190, 77, 213, 249, 31, 111, 200, 94, 57, 113, 179, 238, 205, 77, 231, 108, 202, 139, 116, 222, 65, 210, 210, 66, 140, 42, 129, 183, 7, 71, 198, 238, 42, 250, 173, 229, 36, 250, 54, 213, 128, 46, 241, 251, 67, 20, 209, 115, 102, 28, 106, 62, 195, 124, 14, 218, 180, 216, 179, 182, 139, 160, 253, 4, 45, 85, 75, 117, 32, 38, 158, 247, 166, 230, 145, 201, 119, 228, 246, 204, 190, 63, 25, 189, 245, 135, 100, 154, 166, 73, 0, 103, 74, 157, 179, 79, 49, 173, 87, 204, 118, 4, 218, 125, 238, 6, 24, 223, 95, 208, 243, 226, 115, 235, 78, 217, 108, 243, 179, 34, 73, 123, 238, 110, 7, 175, 47, 193, 246, 26, 50, 3, 230, 145, 93, 93, 97, 141, 217, 236, 47, 198, 39, 219, 55, 59, 11, 28, 37, 154, 238, 168, 23, 90, 167, 244, 58, 231, 28, 186, 2, 254, 15, 0, 0, 255, 255, 128, 76, 233, 191}
	c, err := resources.DecompressContent(b)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if string(c) != string(LoremIpsum) {
		t.Error("c is not equal to LoremIpsum, got:", string(c))
	}
}
