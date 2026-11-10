import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [0, -4, 1, 0] }, { duration: τmax });

  const bars = [...el.querySelectorAll(":scope > path")].reverse();
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
