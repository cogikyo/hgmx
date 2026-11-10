import { animate, τmax } from "../motion.js";

export default function (el) {
  const binoculars = el.querySelector("g");
  animate(binoculars, { translateX: [0, -1.5, 1.5, 0], rotate: [0, -6, 6, 0] }, { duration: τmax });
}
