import { animate, τmax } from "../motion.js";

export default function (el) {
  const water = el.querySelector("path[transform-origin]");
  animate(water, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
}
