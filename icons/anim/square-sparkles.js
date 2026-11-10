import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M21 13"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const spark = el.querySelector('path[d^="M15.41"]');
  animate(spark, { scale: [0.5, 0.75, 1.1, 1], rotate: [0, 70, 90, 90] }, { duration: τmax });

  const glint = el.querySelectorAll('path[d="M11 15H7"], path[d="M9 13v4"]');
  animate(glint, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
