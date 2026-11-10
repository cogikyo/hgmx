import { animate, τmax } from "../motion.js";

export default function (el) {
  const shark = el.querySelector("g");
  animate(shark, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1, 0.3, 0] }, { duration: τmax });
}
