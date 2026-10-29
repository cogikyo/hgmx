import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const petals = el.querySelector('path[d^="M12 16.5A"]');
  animate(petals, { scale: [0, 1.15, 1], rotate: [0, 180] }, { duration: τmax });

  const stamens = ["M12 7.5", "M14.12 9.88", "M16.5", "M14.12 14.12", "M12 16.5V", "m8 16", "M7.5", "m8 8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(stamens, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const heart = el.querySelector("circle");
  animate(heart, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
