import { animate, τmax } from "../motion.js";

export default function (el) {
  const popsicle = el.querySelector("g");
  animate(popsicle, { rotate: [0, -10, 6, 0] }, { duration: τmax });
}
