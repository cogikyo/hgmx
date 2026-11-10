import { animate, τmax } from "../motion.js";

export default function (el) {
  const popsicle = el.querySelector("g");
  animate(popsicle, { translateX: [0, -1, 1, -1, 1, 0] }, { duration: τmax });
}
