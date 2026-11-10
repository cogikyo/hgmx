import { animate, τmax } from "../motion.js";

const turns = new WeakMap();

export default function (el) {
  const Θ = turns.get(el) ?? 0;
  const Δ = Θ + 360 + Math.random() * 100;
  turns.set(el, Δ);

  const moons = el.querySelector("g");
  animate(moons, { rotate: [Θ, Δ] }, { duration: τmax });

  const planet = el.querySelector('circle[r="3.5"]');
  animate(planet, { opacity: [0, 1], scale: [1, 0.5, 1] }, { duration: τmax });
}
