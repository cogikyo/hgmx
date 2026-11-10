import { animate, τmax } from "../motion.js";

export default function (el) {
  const handset = el.querySelector("path");
  animate(handset, { rotate: [0, -15, 12, -8, 0] }, { duration: τmax });
}
