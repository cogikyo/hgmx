import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const skin = el.querySelector("path");
  animate(skin, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pit = el.querySelector("circle");
  animate(pit, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
