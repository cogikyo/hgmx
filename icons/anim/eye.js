import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const lid = el.querySelector("path");
  animate(lid, { scaleY: [1, 0.1, 1] }, { duration: τmid });

  const pupil = el.querySelector("circle");
  animate(pupil, { scale: [1, 0, 1.25, 1] }, { duration: τmax });
}
