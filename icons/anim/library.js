import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const books = ["M4 4", "M8 8", "M12 6"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(books, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const leaning = el.querySelector('path[d^="m16 6"]');
  animate(leaning, { rotate: [-16, 4, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
