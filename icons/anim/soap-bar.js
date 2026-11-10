import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bar = el.querySelector("g");
  animate(bar, { scale: [1, 0.9, 1.05, 1] }, { duration: τmid });

  const bubbles = el.querySelectorAll('circle, path[d="M22 22h.01"]');
  animate(bubbles, { scale: [1, 1.75, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
