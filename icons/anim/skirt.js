import { animate, τmax } from "../motion.js";

export default function (el) {
  const skirt = el.querySelector("g");
  animate(skirt, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
