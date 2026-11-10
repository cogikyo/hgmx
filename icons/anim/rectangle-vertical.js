import { animate, τmax } from "../motion.js";

export default function (el) {
  const shape = el.querySelector("rect");
  animate(shape, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
