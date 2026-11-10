import { animate, τmax } from "../motion.js";

export default function (el) {
  const knife = el.querySelector("path");
  animate(knife, { translateX: [0, -2, 1, 0], translateY: [0, 2, -1, 0] }, { duration: τmax });
}
