import { animate, τmax } from "../motion.js";

export default function (el) {
  const shoe = el.querySelector("g");
  animate(shoe, { rotate: [0, 12, -3, 0] }, { duration: τmax });
}
