-- 2021-10-05 03:06:13.667663+00

CREATE TABLE IF NOT EXISTS "paste" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "content" TEXT NOT NULL,
    "created" TIMESTAMP NOT NULL
);

INSERT INTO "paste" ("id", "content", "created") VALUES
('3E1dp-c5V4DQMDEFbfxr3',	'But I must explain to you how all this mistaken idea of denouncing pleasure and praising pain was born and I will give you a complete account of the system, and expound the actual teachings of the great explorer of the truth, the master-builder of human happiness. No one rejects, dislikes, or avoids pleasure itself, because it is pleasure, but because those who do not know how to pursue pleasure rationally encounter consequences that are extremely painful. Nor again is there anyone who loves or pursues or desires to obtain pain of itself, because it is pain, but because occasionally circumstances occur in which toil and pain can procure him some great pleasure. To take a trivial example, which of us ever undertakes laborious physical exercise, except to obtain some advantage from it? But who has any right to find fault with a man who chooses to enjoy a pleasure that has no annoying consequences, or one who avoids a pain that produces no resultant pleasure? On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee',	'2021-10-04 00:32:07'),
('fouz1Irw2oMnQHoq4nGSk',	'One morning, when Gregor Samsa woke from troubled dreams, he found himself transformed in his bed into a horrible vermin. He lay on his armour-like back, and if he lifted his head a little he could see his brown belly, slightly domed and divided by arches into stiff sections. The bedding was hardly able to cover it and seemed ready to slide off any moment. His many legs, pitifully thin compared with the size of the rest of him, waved about helplessly as he looked. What''s happened to me? he thought. It wasn''t a dream. His room, a proper human room although a little too small, lay peacefully between its four familiar walls. A collection of textile samples lay spread out on the table - Samsa was a travelling salesman - and above it there hung a picture that he had recently cut out of an illustrated magazine and housed in a nice, gilded frame. It showed a lady fitted out with a fur hat and fur boa who sat upright, raising a heavy fur muff that covered the whole of her lower arm towards the viewer. Gregor then turned to look out the window at the dull weather. Drops',	'2021-10-05 00:33:07'),
('Ug1iHQJbu4ZLYyzj5tYiV',	'A wonderful serenity has taken possession of my entire soul, like these sweet mornings of spring which I enjoy with my whole heart. I am alone, and feel the charm of existence in this spot, which was created for the bliss of souls like mine. I am so happy, my dear friend, so absorbed in the exquisite sense of mere tranquil existence, that I neglect my talents. I should be incapable of drawing a single stroke at the present moment; and yet I feel that I never was a greater artist than now. When, while the lovely valley teems with vapour around me, and the meridian sun strikes the upper surface of the impenetrable foliage of my trees, and but a few stray gleams steal into the inner sanctuary, I throw myself down among the tall grass by the trickling stream; and, as I lie close to the earth, a thousand unknown plants are noticed by me: when I hear the buzz of the little world among the stalks, and grow familiar with the countless indescribable forms of the insects and flies, then I feel the presence of the Almighty, who formed us in his own image, and the breath',	'2021-10-05 00:34:07'),
('PePQyBnOsnIJjf7Wsm1Zq',	'{
	"chatID": "2318301",
	"command": "test",
	"message": "test!"
}',	'2021-10-05 09:36:26'),
('ItUMXJYOLbAERrwbUafvt',	'import options from "commander"
import fs from "fs-extra"
import * as env from "./buildSrc/env.js"
import {renderHtml} from "./buildSrc/LaunchHtml.js"
import {spawn, spawnSync} from "child_process"
import path, {dirname} from "path"
import {fetchDictionaries} from ''./buildSrc/DictionaryFetcher.js''
import os from "os"
import {rollup} from "rollup"
import {babelPlugins, bundleDependencyCheckPlugin, getChunkName, resolveLibs} from "./buildSrc/RollupConfig.js"
import {terser} from "rollup-plugin-terser"
import pluginBabel from "@rollup/plugin-babel"
import commonjs from "@rollup/plugin-commonjs"
import {fileURLToPath} from "url"
import nodeResolve from "@rollup/plugin-node-resolve"
import {visualizer} from "rollup-plugin-visualizer"
import glob from "glob"',	'2021-10-05 09:52:59'),
('qGphy2mruYr56S-voJIbK',	'const MINIFY = options.disableMinify !== true
if (!MINIFY) {
	console.warn("Minification is disabled")
}

doBuild().catch(e => {
	console.error(e)
	process.exit(1)
})',	'2021-10-05 09:56:05'),
('UQMSN9hBS6UyI5aNWGpUE',	'<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Decimal Clock</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <h2>Decimal Clock</h2>
    <div id="clock-face"></div>
    <h1 id="clock-time"></h1>
    <script src="script.js"></script>
</body>
</html>',	'2021-10-03 00:30:07'),
('MGgDCnMd9u2x1TtBiJnaX',	'{
  "sdk": {
    "version": "6.0.100-rc.1.21458.32"
  }
}',	'2021-10-04 00:31:07'),
('lcURorZsnQiR9TEy1twvj',	'Far far away, behind the word mountains, far from the countries Vokalia and Consonantia, there live the blind texts. Separated they live in Bookmarksgrove right at the coast of the Semantics, a large language ocean. A small river named Duden flows by their place and supplies it with the necessary regelialia.

It is a paradisematic country, in which roasted parts of sentences fly into your mouth. Even the all-powerful Pointing has no control about the blind texts it is an almost unorthographic life One day however a small line of blind text by the name of Lorem Ipsum decided to leave for the far World of Grammar.

The Big Oxmox advised her not to do so, because there were thousands of bad Commas, wild Question Marks and devious Semikoli, but the Little Blind Text didn’t listen. She packed her seven versalia, put her initial into the belt and made herself on the way. When she reached the first hills of the Italic Mountains, she had a last view back on the skyline of her hometown Bookmarksgrove, the headline of Alphabet Village and the subline of her own road, the Line Lane. Pityful a rethoric question ran over her cheek, then',	'2021-10-05 10:04:14');
