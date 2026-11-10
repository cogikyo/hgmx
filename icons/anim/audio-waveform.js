import { animate, τmax } from "../motion.js";

export default function (el) {
  const wave = el.querySelector("path");
  animate(wave, { scaleY: [1, 0.6, 1.2, 1] }, { duration: τmax });
}
