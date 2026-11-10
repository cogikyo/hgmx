import { animate, τmax } from "../motion.js";

export default function (el) {
  const burger = el.querySelector("g");
  animate(burger, { translateY: [0, -2, 0, 0], scaleY: [1, 1.04, 0.92, 1] }, { duration: τmax });
}
