import { animate, τmax } from "../motion.js";

export default function (el) {
  const shovel = el.querySelector("g");
  animate(shovel, { rotate: [0, -12, 6, 0] }, { duration: τmax });
}
