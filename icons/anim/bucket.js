import { animate, τmax } from "../motion.js";

export default function (el) {
  const bucket = el.querySelector("g");
  animate(bucket, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
