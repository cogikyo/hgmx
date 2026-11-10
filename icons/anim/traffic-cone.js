import { animate, τmax } from "../motion.js";

export default function (el) {
  const cone = el.querySelector("g");
  animate(cone, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
