import { animate, τmax } from "../motion.js";

export default function (el) {
  const skirt = el.querySelector("g");
  animate(skirt, { rotate: [0, -10, 6, -3, 0], scaleX: [1, 1.15, 0.97, 1] }, { duration: τmax });
}
