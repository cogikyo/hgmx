import { animate, τmax } from "../motion.js";

export default function (el) {
  const ribbon = el.querySelector("g");
  animate(ribbon, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
