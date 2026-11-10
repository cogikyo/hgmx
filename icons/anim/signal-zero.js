import { animate } from "../motion.js";

export default function (el) {
  const dot = el.querySelector("path");
  animate(dot, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3 });
}
