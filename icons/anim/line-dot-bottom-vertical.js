import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const dot = el.querySelector("circle");
  animate(dot, { scale: [1, 1.25, 1] }, { duration: τmin });

  const line = el.querySelector("path");
  animate(line, { scaleY: [0, 1] }, { duration: τmid, delay: τmin });
}
