import { animate, τmax } from "../motion.js";

export default function (el) {
  const pickaxe = el.querySelector("g");
  animate(pickaxe, { rotate: [0, -30, 10, 0] }, { duration: τmax });
}
