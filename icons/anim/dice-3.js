import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const die = el.querySelector("g");
  animate(die, { translateY: [0, -3, 1, 0], rotate: [0, -15, 6, 0] }, { duration: τmax });
  const pip = (d) => el.querySelector(`path[d="${d}h.01"]`);
  const diagonal = [pip("M16 8"), pip("M12 12"), pip("M8 16")];
  animate(diagonal, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
