import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const spark = el.querySelector('path[d="m22 2-1.5 1.5"]');
  animate(spark, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3 });
  const bomb = el.querySelector("g");
  animate(bomb, { scale: [1, 0.9, 1.1, 1] }, { duration: τmid, delay: τmin });
}
