import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const stack = el.querySelector("path");
  animate(stack, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const block = el.querySelector("rect");
  animate(block, { translateX: [4, -1, 0], translateY: [-4, 1, 0], rotate: [20, -5, 0], opacity: [0, 1] }, { duration: τmax });
}
