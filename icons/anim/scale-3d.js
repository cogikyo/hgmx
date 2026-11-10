import { animate, τmax } from "../motion.js";

export default function (el) {
  const axes = el.querySelector("g");
  animate(axes, { scale: [1, 0.8, 1.1, 1] }, { duration: τmax });
}
