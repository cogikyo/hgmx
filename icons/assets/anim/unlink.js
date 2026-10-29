import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const upper = el.querySelector('path[d^="m18.84"]');
  animate(upper, { translateX: [0, 2, 0], translateY: [0, -2, 0] }, { duration: τmid });

  const lower = el.querySelector('path[d^="m5.17"]');
  animate(lower, { translateX: [0, -2, 0], translateY: [0, 2, 0] }, { duration: τmid });

  const sparks = el.querySelectorAll("line");
  animate(sparks, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
