import { animate, τmax } from "../motion.js";

export default function (el) {
  const carton = el.querySelector("g");
  animate(carton, { rotate: [0, 12, -4, 0] }, { duration: τmax });
}
