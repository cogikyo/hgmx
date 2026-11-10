import { animate, τmax } from "../motion.js";

export default function (el) {
  const spark = el.querySelector("path");
  animate(spark, { scale: [0.5, 1.25, 1], rotate: [0, 90] }, { duration: τmax });
}
