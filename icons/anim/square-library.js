import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const books = el.querySelectorAll('path[d$="v10"]');
  animate(books, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const leaning = el.querySelector('path[d^="m15 7"]');
  animate(leaning, { rotate: [-16, 4, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
