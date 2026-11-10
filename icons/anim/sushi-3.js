import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const nigiri = el.querySelectorAll('path[d^="M22 10"], path[d^="M20 12"]');
  animate(nigiri, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stripes = ["m6 11", "m10 10", "m14 10", "m18 11"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(stripes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
