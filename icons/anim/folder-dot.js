import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector("path");
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
