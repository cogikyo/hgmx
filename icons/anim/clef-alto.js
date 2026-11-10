import { animate, τmax } from "../motion.js";

export default function (el) {
  const clef = el.querySelector("g");
  animate(clef, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}
