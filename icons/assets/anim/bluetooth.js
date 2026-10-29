import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const rune = el.querySelector("path");
  animate(rune, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(rune, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
