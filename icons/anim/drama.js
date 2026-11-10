import { animate, τmax } from "../motion.js";

export default function (el) {
  const masks = el.querySelector("g");
  animate(masks, { rotate: [0, -8, 5, 0] }, { duration: τmax });
}
