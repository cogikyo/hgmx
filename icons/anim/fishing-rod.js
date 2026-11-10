import { animate, τmax } from "../motion.js";

export default function (el) {
  const rod = el.querySelector("g");
  animate(rod, { rotate: [0, -6, 2, 0] }, { duration: τmax });
}
