import { animate, τmax } from "../motion.js";

export default function (el) {
  const arms = el.querySelector("g");
  animate(arms, { rotate: [0, -45, 10, 0], scale: [1, 0.9, 1] }, { duration: τmax });
}
