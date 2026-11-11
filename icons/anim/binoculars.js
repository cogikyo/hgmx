import { animate, τmax } from "../motion.js";

export default function (el) {
  const binoculars = el.querySelector("g");
  animate(binoculars, { translateX: [0, -3, 2, 0], rotate: [0, -10, 6, 0] }, { duration: τmax });
}
