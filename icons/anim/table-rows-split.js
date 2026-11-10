import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cut = ["M2 10h2", "M8 10h2", "M14 10h2", "M20 10h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(cut, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const top = el.querySelectorAll('path[d^="M3 2v"], path[d="M9 2v4"], path[d="M15 2v4"]');
  animate(top, { translateY: [0, -1.5, 0] }, { duration: τmid, delay: τmin });

  const bottom = el.querySelectorAll('path[d^="M3 22v"], path[d="M3 19h18"], path[d="M9 22v-8"], path[d="M15 22v-8"]');
  animate(bottom, { translateY: [0, 1.5, 0] }, { duration: τmid, delay: τmin });
}
