import { animate, τmax } from "../motion.js";

export default function (el) {
  const sword = el.querySelector("g");
  animate(sword, { rotate: [0, -20, 6, 0] }, { duration: τmax });
}
