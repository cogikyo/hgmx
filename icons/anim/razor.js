import { animate, τmax } from "../motion.js";

export default function (el) {
  const razor = el.querySelector("g");
  animate(razor, { translateX: [0, -2, 1, 0], translateY: [0, 2, -1, 0] }, { duration: τmax });
}
