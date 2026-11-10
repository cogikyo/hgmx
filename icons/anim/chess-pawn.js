import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pawn = el.querySelector("g");
  animate(pawn, { translateY: [0, -3, 1, 0] }, { duration: τmax });
  const base = el.querySelector('path[d^="M5 20"]');
  animate(base, { scaleX: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
