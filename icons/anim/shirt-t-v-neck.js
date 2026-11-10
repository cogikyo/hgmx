import { animate, τmax } from "../motion.js";

export default function (el) {
  const shirt = el.querySelector("g");
  animate(shirt, { scaleY: [1, 1.12, 0.96, 1] }, { duration: τmax });
}
