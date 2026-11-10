import { animate, τmax } from "../motion.js";

export default function (el) {
  const wheel = el.querySelector("g");
  animate(wheel, { rotate: [0, -15, 90] }, { duration: τmax });
}
