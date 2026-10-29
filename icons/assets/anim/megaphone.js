import { animate, τmax } from "../motion.js";

export default function (el) {
  const horn = el.querySelector("g");
  animate(horn, { scale: [1, 0.85, 1.15, 1], rotate: [0, -8, 0] }, { duration: τmax });
}
