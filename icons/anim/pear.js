import { animate, τmax } from "../motion.js";

export default function (el) {
  const pear = el.querySelector("g");
  animate(pear, { rotate: [0, -7, 4, -2, 0] }, { duration: τmax });
}
