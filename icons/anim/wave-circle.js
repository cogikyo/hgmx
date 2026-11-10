import { animate, τmax } from "../motion.js";

export default function (el) {
  const wave = el.querySelector("path");
  animate(wave, { rotate: [0, 35, -10, 0] }, { duration: τmax });
}
