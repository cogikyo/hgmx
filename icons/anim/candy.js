import { animate, τmax } from "../motion.js";

export default function (el) {
  const candy = el.querySelector("g");
  animate(candy, { rotate: [0, 180], scale: [1, 0.85, 1] }, { duration: τmax });
}
