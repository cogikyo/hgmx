import { animate, τmax } from "../motion.js";

export default function (el) {
  const king = el.querySelector("g");
  animate(king, { scaleY: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
